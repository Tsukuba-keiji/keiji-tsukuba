$(function(){
  let param = decodeURI(location.search.slice(1));
  let obj = {};
  param.split('&').forEach(function(param) {
    let queryArr = param.split('=');
    obj[queryArr[0]] = queryArr[1];
  });
  let gr = obj["grade"];
  let cl = obj["class"];
  let text = obj["text"];
  $("#verify").html(`内容：${text} <br> タグ：${gr}年生、${cl}組`);
  
  document.getElementById("button").addEventListener("click",function(){
    obj["crossDomain"]=true;
    let data = JSON.stringify(obj);
    let request = new XMLHttpRequest();
    request.open("get","https://script.google.com/macros/s/AKfycbyPF6-4wT6dljPaT8SPT1LJZP-mPdgwIL6adpC-UMSqwbMskNz0vCZChrJ417PD02TAyA/exec");
    request.setRequestHeader('Content-Type', 'application/json');
    request.send(data);
    request.onreadystatechange = function() {
      if (request.readyState == 4 && request.status == 200) {
        $("body").html("<h1>success<h1>");
      }
    }
  });
});
