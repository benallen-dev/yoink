# YOINK

Yoink. Takes things, quickly.

## Usage

```
go run cmd/main.go
```

## Web UI

The web UI runs on port 8081. You can access it by visiting `http://localhost:8081` in your browser.

The purpose of the UI is to categorise the images that have been yoinked. You view the image and press the appropriate key to move the file to a sorted folder.

Currently there's 4 categories:
- Keep
- Discard
- Nsfw
- Nswf Anime

The reason for these categories is that I'm using this to sort images into training data. Ultimately I want to train a classification model that can detect images I want to keep automatically.

The seperate nsfw and nsfw anime categories is because at least to my mind there's a big difference between photographs and drawn images and I wonder if when training the model that might cause errors. They can always be merged if that doesn't turn out to be true.

