/* eslint-env browser */
const log = msg => {
  document.getElementById('logs').innerHTML += msg + '<br>'
}

window.createSession = async isPublisher => {
  let pc = new RTCPeerConnection({
    iceServers: [
      {
        urls: 'stun:stun.l.google.com:19302'
      }
    ]
  })
  pc.oniceconnectionstatechange = e => log(pc.iceConnectionState)
  pc.onicecandidate = event => {
    if (event.candidate === null) {
      document.getElementById('localSessionDescription').value = btoa(JSON.stringify(pc.localDescription))
    }
  }

  try {
    if (isPublisher) {
      const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: false })
      pc.addStream(document.getElementById('video1').srcObject = stream)
      const offer = await pc.createOffer()
      await pc.setLocalDescription(offer)
    } else {
      pc.addTransceiver('video', {'direction': 'recvonly'})
      const offer = await pc.createOffer();
      await pc.setLocalDescription(offer)
  
      pc.ontrack = event => {
        var el = document.getElementById('video1')
        el.srcObject = event.streams[0]
        el.autoplay = true
        el.controls = true
      }
    }
  }
  catch(err) {
    log(err);
  }
  

  window.startSession = async () => {
    const sdp = document.getElementById('localSessionDescription').value
    const { data:response } = await axios.post('/sdp', sdp);
    document.getElementById('remoteSessionDescription').value = response;

    let sd = document.getElementById('remoteSessionDescription').value
    if (sd === '') {
      return alert('Session Description must not be empty')
    }

    try {
      pc.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(sd))))
    } catch (e) {
      alert(e)
    }
  }

  let btns = document.getElementsByClassName('createSessionButton')
  for (let i = 0; i < btns.length; i++) {
    btns[i].style = 'display: none'
  }

  document.getElementById('signalingContainer').style = 'display: block'
}
