package ledgrid

// Dient dazu, ein Live-Bild ab einer beliebigen, aber ansprechbaren Kamera
// auf dem LED-Grid darzustellen. Als erstes eine Implementation mit Hilfe
// der Video4Linux-Umgebung... nachdem zuerst mal ein paar Konstanten die
// Konfiguration vereinfachen sollen.
// Die 2 möglichen Implementationen der Kamera sind in separaten Dateien
// zu finden, welche über Build-Flags aktiviert werden können:
//
//	-tags=cameraOpenCV
//	-tags=cameraV4L2
//
// Die allgemeinen Konstanten sind:
const (
	camDevName    = "/dev/video0"
	camDevId      = 0
	camWidth      = 320
	camHeight     = 240
	camFrameRate  = 30
	camBufferSize = 4
)

