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
    let data = JSON.stringify(obj);
    let request = new XMLHttpRequest();
    request.open("post","https://script.google.com/macros/s/AKfycbwu2sfN51Dw9OOzudZKjFJu7UF9vzJg8bhhKlY8zaCuJdT7NG_mjwSUlx9yfVisn8ru/exec");
    request.setRequestHeader('Content-Type', 'application/json');
    request.send(data);
  });
});