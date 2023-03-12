# Reverbify

## Overview 
Obsessed with slowed and reverb music on youtube? How about bass boosted tunes and nightcore remixes? Reverbify is application that automates audio processing and applies reverb, bass, and pitch alteration to music videos on youtube. This repository houses the api of the app that creates the infrastructure for building and storing music. 

## Endpoints
### Routes

#### GET https://reverbify-backend-klfqvexjrq-vp.a.run.app/health-check
- Returns back a 200 response with a string ensuring the app is running properly. 

#### POST https://reverbify-backend-klfqvexjrq-vp.a.run.app/signed-url
- Given a file name of a mp3 file on AWS, returns a signed url to that mp3
- Body of the request: 
```yaml
{
   "filename": <place file name here>
}
```

- Response returned
```yaml
{
   "signedUrl": <Link to file on AWS S3>
}
```

#### POST https://reverbify-backend-klfqvexjrq-vp.a.run.app/reverb-song
- Given a youtube link and a set of audio processing requests, returns meta data about the produced song along with a link to file on AWS
- Body of the request: 
```yaml
{
    "url": "<youtube link>",
    "pitch": "<double from .5 to 1.5 wrapped in quotes, 1 for original form>",
    "reverb": "<number from 1 through 11 represents a type of reverb, 0 for no reverb>",
    "bass": {
        "change": <boolean true if bass, false if not>, 
        "centerFreq": "60",
        "filterWidth": "50",
        "gain": "<integer from 5-12 indicating how loud the bass is in decibals>"
    }
}
```
- reverb types: 
```golang
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
  ```

- Response returned
```yaml
{
    "title": <title of youtube video>,
    "author": <channel that published video>,
    "duration": <time of newly created mp3 file>,
    "thumbnail": <link to thumbnail of youtube video>,
    "signedUrl": <link to file on AWS S3>,
    "filename": <filename of song on AWS>,
    "timestamp": <date of song created>
}
```




## Tech Stack
### Languages: 
- Golang (Gin): Used the gin framework to create a REST API that handles user requests for processing music and retrieving finished products on AWS S3

### Devops: 
- AWS S3: cloud storage used to house all the music creations for each user
- Docker: used to containerize the app so that each clone of the application can be given the same access to the required executable programs
- Google Cloud Platforms: leveraged the container registry for hosting an image of the api such that its accessible to the mobile app

### Services: 
- FFMPEG: command line executable program used for audio processing (bass bossting, pitch alteration, applying reverb)
- yt-dlp: command line executable program used to download music from youtube given a valid url to a music 
- Youtube API: service provided by youtube to provide meta information about a video given a video id

