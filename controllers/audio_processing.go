package controllers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func Init_audio_processing(c *gin.Context) {
	var body audio_request
	err := c.BindJSON(&body)

	if handleError(err, c, "Invalid request body") {
		return
	}

	transform(c, body.Url, body.Pitch, body.Reverb, body.Bass)
}

func Health_check(c *gin.Context) {
	healthCheck(c)
}

func transform(c *gin.Context, url string, pitch string, reverb string, bass bass) {
	var videoId string
	url, videoId = processUrl(url)

	title, fileName, author, err := getTitle(url, videoId)

	if handleError(err, c, title) {
		return
	}

	// download video
	fileName, err = getMP3FromYotube(url, fileName)
	if handleError(err, c, "Failed to download from youtube.") {
		return
	}

	// add reverb
	didReverb, fileNameOutput := processReverb(fileName, c, reverb)
	if !didReverb {
		return
	}

	// change bass
	didBass, fileNameOutput := processBass(fileNameOutput, c, bass)
	if !didBass {
		return
	}

	// change pitch
	didPitch, finalPath := processPitch(fileNameOutput, c, pitch)
	if !didPitch {
		return
	}

	// meta data
	thumbnailURL := getThumbnail(url)
	duration := getVideoLength(finalPath)
	fmt.Println("Complete!")
	signedUrl, err := upload(finalPath, fileNameOutput)
	if handleError(err, c, signedUrl) {
		return
	}
	deleteFile(finalPath)
	sendAudioResponse(c, title, duration, author, thumbnailURL, signedUrl, fileNameOutput)
}

func processUrl(url string) (string, string) {
	regex := regexp.MustCompile(`^.*(youtu.be\/|v\/|u\/\w\/|embed\/|watch\?v=|\&v=)([^#\&\?]*).*`)
	videoId := regex.FindStringSubmatch(url)[2]
	return "https://www.youtube.com/watch?v=" + videoId, videoId
}

func getMP3FromYotube(url string, fileName string) (string, error) {
	fileNameMP4 := fileName + ".mp4"
	fileNameMP3 := fileName + "_" + createTimeStamp() + ".mp3"

	// Uses youtube-dl exec on machine to download videos from youtube
	fmt.Println("Downloaded mp4 file...")
	downloadCommand := exec.Command("yt-dlp", "-f", "ba", "-S", "ext:mp4", "-o", fileNameMP4, url)
	downloadOutput, err := downloadCommand.CombinedOutput()
	if logErr(err, downloadOutput) {
		return "", err
	}

	// converts mp4 to mp3 using ffmpeg
	fmt.Println("Converting mp4 to mp3 file...")
	convertCommand := exec.Command("ffmpeg", "-i", fileNameMP4, fileNameMP3)
	convertOutput, err := convertCommand.CombinedOutput()
	if logErr(err, convertOutput) {
		return "", err
	}

	// removes uneeded mp4 file
	fmt.Println("Removing mp4 file...")
	err = deleteFile(fileNameMP4)
	if err != nil {
		return "", err
	}

	fmt.Println("Successfully downloaded and converted YouTube video to MP3!")
	return fileNameMP3, nil
}

func processReverb(fileNameInput string, c *gin.Context, reverb string) (bool, string) {
	var filter_path string = "./ReverbFilters/"
	reverbTypes := map[string]string{
		"1":  "CUSTOM_pump_verb.WAV",
		"2":  "INSTR_snare_gate.WAV",
		"3":  "VOC_deep_verb.WAV",
		"4":  "VOC_good_ol'_verb.WAV",
		"5":  "VOC_slap_hall.WAV",
		"6":  "VOC_vocal_magic.WAV",
		"7":  "ORCH_small_hall.WAV",
		"8":  "ORCH_medium_hall.WAV",
		"9":  "ORCH_large_hall.WAV",
		"10": "ORCH_concert_hall.WAV",
		"11": "LIVE_live_arena.WAV",
	}

	fileNameOutput := fileNameInput
	if filter, ok := reverbTypes[reverb]; ok {
		filter_path = filter_path + filter
		fmt.Println("Adding reverb...")
		fileNameOutput = "rev_" + reverb + "_" + fileNameInput
		reverbCommand := exec.Command("ffmpeg", "-i", fileNameInput, "-i", filter_path, "-filter_complex",
			"[0] [1] afir=dry=10:wet=10 [reverb]; [0] [reverb] amix=inputs=2:weights=10 1", fileNameOutput)
		output, err := reverbCommand.CombinedOutput()
		if handleError(err, c, "Failed in the reverb process.") || handleError(deleteFile(fileNameInput), c, "failed deleting file.") {
			logErr(err, output)
			return false, ""
		}
	}
	return true, fileNameOutput
}

func processPitch(fileNameInput string, c *gin.Context, pitch string) (bool, string) {
	path := "./music/pitch_speed:_" + pitch + "_" + fileNameInput
	pitchCommand := exec.Command("ffmpeg", "-i", fileNameInput, "-af", "asetrate=44100*"+pitch+",aresample=44100", path)
	output, err := pitchCommand.CombinedOutput()
	fmt.Println("Altering pitch to " + pitch)
	if handleError(err, c, "Failed in the pitch altering process.") || handleError(deleteFile(fileNameInput), c, "failed deleting file.") {
		logErr(err, output)
		return false, ""
	}
	return true, path
}

func processBass(fileNameInput string, c *gin.Context, bass bass) (bool, string) {
	fileNameOutput := fileNameInput
	if bass.SetBass {
		fileNameOutput = "bass_" + fileNameInput
		bassArgs := []string{"-i", fileNameInput, "-af", "equalizer=f=" + bass.CentFreq + ":width_type=h:width=" + bass.FilterWidth + ":g=" + bass.Gain, fileNameOutput}
		bassCommand := exec.Command("ffmpeg", bassArgs...)
		output, err := bassCommand.Output()
		if handleError(err, c, "Failed in the bass process.") || handleError(deleteFile(fileNameInput), c, "failed deleting file.") {
			logErr(err, output)
			return false, ""
		}
	}

	return true, fileNameOutput
}

func getTitle(url string, videoId string) (string, string, string, error) {
	apiKey := os.Getenv("API_KEY")

	service, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		fmt.Println(fmt.Sprint(err))
		return "", "", "", err
	}

	videoResponse, err := service.Videos.List([]string{"snippet"}).Id(videoId).Do()
	if err != nil {
		fmt.Println(fmt.Sprint(err))
		return "", "", "", err
	}

	title := videoResponse.Items[0].Snippet.Title
	publisher := videoResponse.Items[0].Snippet.ChannelTitle

	raw := string(title)
	title_vid := raw
	if len(raw) > 2 {
		title = raw[:len(raw)-1]
	}
	fileName := strings.Replace(title_vid, " ", "_", -1)

	return title, fileName, publisher, err
}

func getThumbnail(url string) string {
	// compile regex expression
	regex := regexp.MustCompile(`v=([^&]+)`)

	// find instance within vid url
	videoID := regex.FindStringSubmatch(url)[1]

	return "https://img.youtube.com/vi/" + videoID + "/sddefault.jpg"
}

func getVideoLength(fileName string) string {
	readCommand := exec.Command("ffmpeg", "-i", fileName)
	readOutput, _ := readCommand.CombinedOutput()

	// search for video duration
	regex := regexp.MustCompile(`Duration:\s([0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{2})`)

	// find instance within vid url
	videoDuration := regex.FindStringSubmatch(string(readOutput))[1]

	return videoDuration
}

func createTimeStamp() string {
	t := time.Now()
	timestamp := t.Format("2006-01-02 15:04:05")
	modifiedTimestamp := strings.Replace(timestamp, " ", "_", -1)

	return modifiedTimestamp
}

func deleteFile(path string) error {
	deleteCommand := exec.Command("rm", "-r", path)
	_, err := deleteCommand.CombinedOutput()
	return err
}

func logErr(err error, output []byte) bool {
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
	}
	return err != nil
}
