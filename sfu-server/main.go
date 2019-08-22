package main

import (
	"io"
	"log"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v2"

	"github.com/sikang99/pion-sfu-example/internal/signal"
)

const (
	// PLI (Pictire Loss Indication)
	rtcpPLIInterval = time.Second * 3
)

func main() {
	// default setting for logger
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	sdpInChan, sdpOutChan := signal.HTTPSDPServer()

	// Everything below is the Pion WebRTC API, thanks for using it ❤️.
	// Create a MediaEngine object to configure the supported codec
	m := webrtc.MediaEngine{}

	//util.DebugPrintf("here")
	// Setup the codecs you want to use.
	// Only support VP8, this makes our proxying code simpler
	m.RegisterCodec(webrtc.NewRTPVP8Codec(webrtc.DefaultPayloadTypeVP8, 90000))
	log.Println("VP8 codec support")
	m.RegisterCodec(webrtc.NewRTPH264Codec(webrtc.DefaultPayloadTypeH264, 90000)) // added by sikang
	log.Println("H264 codec support")

	// Create the API object with the MediaEngine
	api := webrtc.NewAPI(webrtc.WithMediaEngine(m))

	offer := webrtc.SessionDescription{}
	signal.Decode(<-sdpInChan, &offer)
	log.Println("OFFER\n", offer) // json format

	peerConnectionConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := api.NewPeerConnection(peerConnectionConfig)
	if err != nil {
		log.Fatalln(err)
	}

	// Allow us to receive 1 video track
	if _, err = peerConnection.AddTransceiver(webrtc.RTPCodecTypeVideo); err != nil {
		log.Fatalln(err)
	}

	localTrackChan := make(chan *webrtc.Track)
	// Set a handler for when a new remote track starts, this just distributes all our packets
	// to connected peers
	peerConnection.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {
		// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
		// This can be less wasteful by processing incoming RTCP events, then we would emit a NACK/PLI when a viewer requests it
		go func() {
			ticker := time.NewTicker(rtcpPLIInterval)
			for range ticker.C {
				if rtcpSendErr := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: remoteTrack.SSRC()}}); rtcpSendErr != nil {
					log.Println(rtcpSendErr)
				}
			}
		}()

		// Create a local track, all our SFU clients will be fed via this track
		localTrack, newTrackErr := peerConnection.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), "video", "pion")
		if newTrackErr != nil {
			log.Fatalln(newTrackErr)
		}
		localTrackChan <- localTrack

		rtpBuf := make([]byte, 1400)
		for {
			i, readErr := remoteTrack.Read(rtpBuf)
			if readErr != nil {
				log.Fatalln(err)
			}

			// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
			if _, err = localTrack.Write(rtpBuf[:i]); err != nil && err != io.ErrClosedPipe {
				log.Fatalln(err)
			}
		}
	})

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		log.Fatalln(err)
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("ANSWER for Caster\n", answer) // json format of SDP
	// Get the LocalDescription and take it to base64 so we can paste in browser
	baseAnswer := signal.Encode(answer)
	sdpOutChan <- baseAnswer
	log.Println(baseAnswer)

	localTrack := <-localTrackChan
	for {
		log.Println("Now wait to start sendonly peer connection")

		recvOnlyOffer := webrtc.SessionDescription{}
		signal.Decode(<-sdpInChan, &recvOnlyOffer)
		log.Println(recvOnlyOffer) // json format

		// Create a new PeerConnection
		peerConnection, err := api.NewPeerConnection(peerConnectionConfig)
		if err != nil {
			log.Fatalln(err)
		}

		_, err = peerConnection.AddTrack(localTrack)
		if err != nil {
			log.Fatalln(err)
		}

		// Set the remote SessionDescription
		err = peerConnection.SetRemoteDescription(recvOnlyOffer)
		if err != nil {
			log.Fatalln(err)
		}

		// Create answer
		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			log.Println(err)
			panic(err)
		}

		// Sets the LocalDescription, and starts our UDP listeners
		err = peerConnection.SetLocalDescription(answer)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("ANSWER for Player\n", answer) // json format of SDP
		// Get the LocalDescription and take it to base64 so we can paste in browser
		playAnswer := signal.Encode(answer)
		sdpOutChan <- playAnswer
		log.Println(playAnswer)
	}
}
