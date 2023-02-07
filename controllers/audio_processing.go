package controllers

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var filter_path string = "./LexiconPCM90_Halls/CUSTOM_pump_verb.WAV"

func Init_audio_processing(c *gin.Context) {
	var body audio_request
	err := c.BindJSON(&body)

	if handleError(err, c, "Invalid request body") {
		return
	}

	// original pitch
	var pitch string = "1.0"
	if body.PitchType == 1 {
		//  1 fast (nightcore)
		pitch = "1.15"
	} else if body.PitchType == -1 {
		// -1 daycore (slow)
		pitch = "0.85"
	}

	transform(c, body.Url, filter_path, pitch)
}

func transform(c *gin.Context, url string, filter string, pitch string) {

	url = processUrl(url)

	title, fileName, err := getTitle(url)
	if handleError(err, c, title) {
		return
	}

	// download video
	err = getMP3FromYotube(url, fileName)
	if handleError(err, c, "Failed to download from youtube.") {
		return
	}

	// add reverb
	fmt.Println("Adding reverb...")
	fileNameInput := fileName + ".mp3"
	fileNameOutput := "reverb_" + fileNameInput
	reverbCommand := exec.Command("ffmpeg", "-i", fileNameInput, "-i", filter, "-filter_complex",
		"[0] [1] afir=dry=10:wet=10 [reverb]; [0] [reverb] amix=inputs=2:weights=10 1", fileNameOutput)
	_, err = reverbCommand.CombinedOutput()
	if handleError(err, c, "Failed in the reverb process.") || handleError(deleteFile(fileNameInput), c, "failed deleting file.") {
		return
	}

	// alter pitch
	var core string = "norm"
	if pitch == "1.15" {
		core = "fast"
	} else if pitch == "0.85" {
		core = "slow"
	}

	path := "./music/pitch_" + core + "_" + fileNameOutput
	fmt.Println("Lowering pitch...")
	pitchCommand := exec.Command("ffmpeg", "-i", fileNameOutput, "-af", "asetrate=44100*"+pitch+",aresample=44100", path)
	_, err = pitchCommand.CombinedOutput()
	if handleError(err, c, "Failed in the pitch altering process.") || handleError(deleteFile(fileNameOutput), c, "failed deleting file.") {
		return
	}

	thumbnailURL := getThumbnail(url)
	duration := getVideoLength(path)
	fmt.Println("Complete!")
	upload(path, fileNameInput)
	sendAudioResponse(c, title, duration, thumbnailURL)
}

func processUrl(url string) string {
	regex := regexp.MustCompile(`^.*(youtu.be\/|v\/|u\/\w\/|embed\/|watch\?v=|\&v=)([^#\&\?]*).*`)
	videoID := regex.FindStringSubmatch(url)[2]
	return "https://www.youtube.com/watch?v=" + videoID
}

func getMP3FromYotube(url string, fileName string) error {
	fileNameMP4 := fileName + ".mp4"
	fileNameMP3 := fileName + ".mp3"

	// Uses youtube-dl exec on machine to download videos from youtube
	fmt.Println("Downloaded mp4 file...")
	downloadCommand := exec.Command("yt-dlp", "-f", "ba", "-S", "ext:mp4", "-o", fileNameMP4, url)
	downloadOutput, err := downloadCommand.CombinedOutput()
	if logErr(err, downloadOutput) {
		return err
	}

	// converts mp4 to mp3 using ffmpeg
	fmt.Println("Converting mp4 to mp3 file...")
	convertCommand := exec.Command("ffmpeg", "-i", fileNameMP4, fileNameMP3)
	convertOutput, err := convertCommand.CombinedOutput()
	if logErr(err, convertOutput) {
		return err
	}

	// removes uneeded mp4 file
	fmt.Println("Removing mp4 file...")
	err = deleteFile(fileNameMP4)
	if err != nil {
		return err
	}

	fmt.Println("Successfully downloaded and converted YouTube video to MP3!")
	return nil
}

func getTitle(url string) (string, string, error) {
	getTitleCommand := exec.Command("youtube-dl", "--skip-download", "--get-title", url)
	getTitleOutput, err := getTitleCommand.CombinedOutput()
	if logErr(err, getTitleOutput) {
		return "ERROR: unable to get title.", "", err
	}

	raw := string(getTitleOutput)
	title := raw
	if len(raw) > 2 {
		title = raw[:len(raw)-1]
	}
	fileName := strings.Replace(title, " ", "_", -1)

	return title, fileName, nil
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
