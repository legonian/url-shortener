const clickCount = document.getElementById('number_of_clicks')
const userInput = document.getElementById('url_text')
const copyButton = document.getElementById('copy_button')

copyButton.onclick = function () {
  userInput.select()
  navigator.clipboard.writeText(userInput.value)
}

; (async () => {
  const hostname = window.location.hostname
  const port = window.location.port === '' ? '' : `:${window.location.port}`
  userInput.value = `${hostname}${port}/${userInput.value}`

  userInput.select()
})()
