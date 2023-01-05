package controllers

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// This command is used to update and install the FFmpeg package on a system that uses
// RUN apt-get -y update && apt-get -y upgrade && apt-get install -y --no-install-recommends ffmpeg

var filter_path string = "./LexiconPCM90_Halls/"

type audio_request struct {
	Url       string `json:"url"`
	PitchType int    `json:"pitch"` // 0 nothing, 1 fast (nightcore), -1 daycore (slow)
}

func Test(c *gin.Context) {
	fn := "reverb_dummy.mp3"
	upload("./music/"+fn, fn)
	c.JSON(200, "")
}

func Init_audio_processing(c *gin.Context) {
	// youtube link
	var body audio_request
	if err := c.BindJSON(&body); err != nil {
		c.String(400, "Invalid request body")
		return
	}
	var pitch string = "1.0"
	if body.PitchType == 1 {
		pitch = "1.15"
	} else if body.PitchType == -1 {
		pitch = "0.85"
	}
	transform(c, body.Url, filter_path+"CUSTOM_pump_verb.WAV", pitch)
}

func transform(c *gin.Context, url string, filter string, pitch string) {
	// get title information
	title, fileName := getTitle(url)
	if title == fileName && title == "ERROR: unable to get title." {
		c.JSON(400, gin.H{
			"ERROR": "Failed to get youtube header information.",
		})
		return
	}

	// download video
	if !getMP3FromYotube(url, fileName) {
		c.JSON(400, gin.H{
			"ERROR": "Failed to download from youtube. ",
		})
		return
	}

	// add reverb
	fileName = fileName + ".mp3"
	fileNameRev := "reverb_" + fileName
	fmt.Println("Adding reverb...")
	reverbCommand := exec.Command("ffmpeg", "-i", fileName, "-i", filter, "-filter_complex",
		"[0] [1] afir=dry=10:wet=10 [reverb]; [0] [reverb] amix=inputs=2:weights=10 1", fileNameRev)
	reverbOutput, err := reverbCommand.CombinedOutput()
	if logErr(err, reverbOutput) || !deleteFile(fileName) {
		c.JSON(400, gin.H{
			"ERROR": "Failed in the reverbing process or deleting excess file.",
		})
		return
	}

	// alter pitch
	var core string = "norm"
	if pitch == "1.15" {
		core = "fast"
	} else if pitch == "0.85" {
		core = "slow"
	}
	fileNamePit := "pitch_" + core + "_" + fileNameRev
	fmt.Println("Lowering pitch...")
	path := "./music/" + fileNamePit
	pitchCommand := exec.Command("ffmpeg", "-i", fileNameRev, "-af", "asetrate=44100*"+pitch+",aresample=44100", path)
	pitchOutput, err := pitchCommand.CombinedOutput()
	if logErr(err, pitchOutput) || !deleteFile(fileNameRev) {
		c.JSON(400, gin.H{
			"ERROR": "Failed in the altering pitch process or deleting excess file.",
		})
		return
	}

	thumbnailURL := getThumbnail(url)
	duration := getVideoLength(path)
	fmt.Println("Complete!")

	c.JSON(200, gin.H{
		"title":     title,
		"duration":  duration,
		"thumbnail": thumbnailURL,
	})
}

func getMP3FromYotube(url string, fileName string) bool {
	fileNameMP4 := fileName + ".mp4"
	fileNameMP3 := fileName + ".mp3"

	// Uses youtube-dl exec on machine to download videos from youtube
	fmt.Println("Downloaded mp4 file...")
	downloadCommand := exec.Command("yt-dlp", "-f", "ba", "-S", "ext:mp4", "-o", fileNameMP4, url)
	downloadOutput, err := downloadCommand.CombinedOutput()
	if logErr(err, downloadOutput) {
		return false
	}

	// converts mp4 to mp3 using ffmpeg
	fmt.Println("Converting mp4 to mp3 file...")
	convertCommand := exec.Command("ffmpeg", "-i", fileNameMP4, fileNameMP3)
	convertOutput, err := convertCommand.CombinedOutput()
	if logErr(err, convertOutput) {
		return false
	}

	// removes uneeded mp4 file
	fmt.Println("Removing mp4 file...")
	if !deleteFile(fileNameMP4) {
		return false
	}

	fmt.Println("Successfully downloaded and converted YouTube video to MP3!")
	return true
}

func getTitle(url string) (string, string) {
	getTitleCommand := exec.Command("youtube-dl", "--skip-download", "--get-title", url)
	getTitleOutput, err := getTitleCommand.CombinedOutput()
	if logErr(err, getTitleOutput) {
		return "ERROR: unable to get title.", "ERROR: unable to get title."
	}

	raw := string(getTitleOutput)
	title := raw
	if len(raw) > 2 {
		title = raw[:len(raw)-1]
	}
	fileName := strings.Replace(title, " ", "_", -1)

	return title, fileName
}

func getThumbnail(url string) string {
	// compile regex expression
	regex := regexp.MustCompile(`v=([^&]+)`)

	// find instance within vid url
	videoID := regex.FindStringSubmatch(url)[1]

	return "https://img.youtube.com/vi/" + videoID + "/maxresdefault.jpg"
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

func deleteFile(path string) bool {
	deleteCommand := exec.Command("rm", "-r", path)
	deleteOutput, err := deleteCommand.CombinedOutput()
	return !logErr(err, deleteOutput)
}

func logErr(err error, output []byte) bool {
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
	}
	return err != nil
}
