# Reverbify

## TODO
- Connect to MongoDB
- Figure out BD schema for user auth, users specific songs, and how playlist organization works
- Connect to AWS S3 to store mp3 files
- add golang web framework to create a simple rest API


this is something to consider later when using docker and adding FFmpeg
// This command is used to update and install the FFmpeg package on a system that uses

## Setup and Install
`~/Reverbify go get`
`apt-get -y update && apt-get -y upgrade && apt-get install -y --no-install-recommends ffmpeg`

**Preliminary thoughts / incomplete**

## Schema

User Collection

{
username
password
playlists: [default: {id, song_list: []}, ]
}

Song Collection
{
id: 341321
spedup: 
binary: actual file
}


Playlist Collection
