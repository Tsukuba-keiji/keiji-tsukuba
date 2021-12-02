$(function(){
  let param = decodeURIComponent(location.search.slice(1));
  let obj = {};
  param.split('&').forEach(function(param) {
    let queryArr = param.split('=');
    obj[queryArr[0]] = queryArr[1];
  });
  let dt = {
    "change":"授業変更",
    "thing":"持ち物連絡",
    "other":"その他"
  };
  let gr = obj["grade"];
  let cl = obj["class"];
  let ki = dt[obj["kind"]];
  let text = obj["text"].replace(/\+/g," ");
  $("#verify").html(`内容：${text} <br> タグ：${gr}年生、${cl}組、${ki}`);
  
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
        $("body").html("<h1>アップロードに成功しました！</h1><h3>検索から確認してください</h3>");
      })();
  });
});
