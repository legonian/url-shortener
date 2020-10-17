const clickCount = document.getElementById('number_of_clicks')
const userInput = document.getElementById('url_text')
const copyButton = document.getElementById('copy_button')

async function getInfo(){
  const urlCode = window.location.pathname.split('/')[1]
  const formBody = JSON.stringify({
    url: urlCode
  })

  const res = await window.fetch(`/api/${urlCode}`, {
    method: 'post',
    headers: { 'Content-Type': 'application/json;charset=UTF-8' },
    body: formBody
  })

  if (res.status === 200) {
    return await res.json()
  } else {
    return false
  }
}

copyButton.onclick = function () {
  userInput.select()
  navigator.clipboard.writeText(userInput.value)
}

; (async () => {
  const urlInfo = await getInfo()
  if (urlInfo && urlInfo.ok){
    clickCount.innerText = urlInfo.views_count

    const hostname = window.location.hostname
    const port = window.location.port == 80 || window.location.port == 443 ?
      '' : `:${window.location.port}`
    const pathname = urlInfo.short_url
    const shortURL = `${hostname}${port}/${pathname}`
    userInput.defaultValue = shortURL

    userInput.select()
  }
})()
