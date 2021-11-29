$(function(){
  let request = new XMLHttpRequest();
  request.open("GET","https://script.google.com/macros/s/AKfycbx8z9rPLt_lGPzUblew1IFGpE20yzhHfFqsjCcbhJJV2Gh15dVO1DdMNV5OwGYECOtC/exec");
  let param = location.search.slice(1);
  let obj = {};
  param.split('&').forEach(function(param) {
    let queryArr = param.split('=');
    obj[queryArr[0]] = queryArr[1];
  });
  let gr = obj["grade"];
  let cl = obj["class"];
  request.onreadystatechange = function() {
    if (request.readyState == 4 && request.status == 200) {
      let text = JSON.parse(request.responseText);
      search(text,gr,cl);
    }
  };
  request.send(null);
});

function search(data,gr,cl){
  let result = document.getElementById("result");
  for(let i of data["tags"]){
    console.log(gr);
    if(i["grade"]==gr && i["class"]==cl){
      result.insertAdjacentHTML("beforeend",data["text"][i["pointer"]]+"<br><hr><br>");
    }
  }
}
