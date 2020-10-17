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

; (async () => {
  const urlInfo = await getInfo()
  if (!urlInfo){
    return
  }
  console.log('urlInfo =', urlInfo)
})()