$(function(){
  let param = location.search.slice(1);
  let obj = {};
  param.split('&').forEach(function(param) {
    let queryArr = param.split('=');
    obj[queryArr[0]] = queryArr[1].replace(/+/g," ");
  });
  let gr = decodeURIComponent(obj["grade"]);
  let cl = decodeURIComponent(obj["class"]);
  let text = decodeURIComponent(obj["text"]);
  $("#verify").html(`内容：${text} <br> タグ：${gr}年生、${cl}組`);
  
  document.getElementById("button").addEventListener("click",function(){
    obj["crossDomain"]=true;
    let data = JSON.stringify(obj);
    document.getElementById("button").disabled = true;
     (async () => {
        const resp = await fetch('https://script.google.com/macros/s/AKfycbxopKKU6ECZnC1Ja99QYGAKuxAt4MRAk42BcVkvXS-TiKrSQbeTdE36sflq0Gsx76fG/exec', {
          method: 'POST',
          body: data,
          headers: {
            'Content-Type': 'text/plain'
          }
        });
        console.log(`2: ${await resp.text()}`); // Ok
        $("body").html("<h1>success</h1>");
      })();
  });
});
