package main

import (
	"fmt"
	"os/exec"
)

// This command is used to update and install the FFmpeg package on a system that uses
// RUN apt-get -y update && apt-get -y upgrade && apt-get install -y --no-install-recommends ffmpeg

func main() {
	transform("heheheh", "https://www.youtube.com/watch?v=fe-CdBzr9Kg", "VOC_deep_verb.WAV")
}

func transform(fileName string, url string, filter string) {
	// download video
	getMP3FromYotube(url, fileName)

	// add reverb
	fileName = fileName + ".mp3"
	fileNameRev := "r" + fileName
	fmt.Println("Adding reverb...")
	reverbCommand := exec.Command("ffmpeg", "-i", fileName, "-i", filter, "-filter_complex",
		"[0] [1] afir=dry=10:wet=10 [reverb]; [0] [reverb] amix=inputs=2:weights=10 1", fileNameRev)
	reverbOutput, err := reverbCommand.CombinedOutput()
	logErr(err, reverbOutput)

	// lower pitch
	fileNamePit := "p" + fileNameRev
	fmt.Println("Lowering pitch...")
	pitchCommand := exec.Command("ffmpeg", "-i", fileNameRev, "-af", "asetrate=44100*0.8,aresample=44100", fileNamePit)
	pitchOutput, err := pitchCommand.CombinedOutput()
	logErr(err, pitchOutput)

	fmt.Println("Complete!")

}

func getMP3FromYotube(url string, fileName string) {
	fileNameMP4 := fileName + ".mp4"
	fileNameMP3 := fileName + ".mp3"

	// Uses youtube-dl exec on machine to download videos from youtube
	fmt.Println("Downloaded mp4 file...")
	downloadCommand := exec.Command("youtube-dl", "-f", "best", "-o", fileNameMP4, url)
	downloadOutput, err := downloadCommand.CombinedOutput()
	logErr(err, downloadOutput)

	// converts mp4 to mp3 using ffmpeg
	fmt.Println("Converting mp4 to mp3 file...")
	convertCommand := exec.Command("ffmpeg", "-i", fileNameMP4, fileNameMP3)
	convertOutput, err := convertCommand.CombinedOutput()
	logErr(err, convertOutput)

	// removes uneeded mp4 file
	fmt.Println("Removing mp4 file...")
	deleteFile(fileNameMP4)

	fmt.Println("Successfully downloaded and converted YouTube video to MP3!")
}

func deleteFile(fileName string) {
	deleteCommand := exec.Command("rm", "-r", fileName)
	deleteOutput, err := deleteCommand.CombinedOutput()
	logErr(err, deleteOutput)
}

func logErr(err error, output []byte) {
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return
	}
}
