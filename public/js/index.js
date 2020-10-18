const userInput = document.getElementById('url_text')
const submitButton = document.getElementById('submit_url_button')
const errorMessage = document.getElementById('error_message')

const showError = () => { errorMessage.hidden = false }
const hideError = () => { errorMessage.hidden = true }

async function sendToServer(urlText){
  console.log('urlText =', urlText)
  hideError()

  const formBody = JSON.stringify({
    url: urlText
  })

  const res = await window.fetch('/create', {
    method: 'post',
    headers: { 'Content-Type': 'application/json;charset=UTF-8' },
    body: formBody
  })

  if (res.status === 200) {
    ans = await res.json()
    
    const protocol =  window.location.protocol
    const hostname = window.location.hostname
    const port = window.location.port
    const url = `${protocol}//${hostname}:${port}/${ans.short_url}/info`
    window.location.href = url
  } else {
    showError()
  }
}

submitButton.onclick = function () { sendToServer(userInput.value) }

userInput.addEventListener("keyup", function (event){
  if (event.key === 'Enter') {
    submitButton.click()
  }
})
