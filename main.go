package main

import (
	"fmt"
	"os/exec"
)

// This command is used to update and install the FFmpeg package on a system that uses
// RUN apt-get -y update && apt-get -y upgrade && apt-get install -y --no-install-recommends ffmpeg

// exectuable command to slow down music by 50 percent
// ffmpeg -i input.mp3 -filter:a "atempo=0.5" output.mp3

/*

This command will take an input file called "input.mp3" and create a new output file called "output.mp3"
with a moderate amount of reverb added to the audio. The four parameters of the aecho filter control the
amount of echo, the decay factor, the delay in milliseconds, and the wet/dry mix, respectively.

You can adjust the parameters of the aecho filter to achieve the desired amount of reverb for your MP3 file.
For example, increasing the decay factor and the wet/dry mix will result in a stronger reverb effect, while
decreasing them will result in a weaker effect.

*/
// ffmpeg -i input.mp3 -filter:a "aecho=0.8:0.9:1000:0.3" output.mp3

func main() {
	getMP3FromYotube("https://www.youtube.com/watch?v=ZBRuPESPiog")
}

func getMP3FromYotube(url string) {
	// Uses youtube-dl exec on machine to download videos from youtube
	fmt.Println("Downloaded mp4 file...")
	downloadCommand := exec.Command("youtube-dl", "-f", "best", "-o", "video.mp4", url)
	downloadOutput, err := downloadCommand.CombinedOutput()
	logErr(err, downloadOutput)

	// converts mp4 to mp3 using ffmpeg
	fmt.Println("Converting mp4 to mp3 file...")
	convertCommand := exec.Command("ffmpeg", "-i", "video.mp4", "video.mp3")
	convertOutput, err := convertCommand.CombinedOutput()
	logErr(err, convertOutput)

	// removes uneeded mp4 file
	fmt.Println("Removing mp4 file...")
	deleteFile("video.mp4")

	fmt.Println("Successfully downloaded and converted YouTube video to MP3!")
}

func deleteFile(fileName string) {
	deleteCommand := exec.Command("rm", "-r", "video.mp4")
	deleteOutput, err := deleteCommand.CombinedOutput()
	logErr(err, deleteOutput)
}

func logErr(err error, output []byte) {
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return
	}
}
