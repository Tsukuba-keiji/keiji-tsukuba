$(function(){
  let request = new XMLHttpRequest();
  request.open("GET","https://script.google.com/macros/s/AKfycbwu2sfN51Dw9OOzudZKjFJu7UF9vzJg8bhhKlY8zaCuJdT7NG_mjwSUlx9yfVisn8ru/exec");
  let param = location.search.slice(1);
  let obj = {};
  param.split('&').forEach(function(param) {
    let queryArr = param.split('=');
    obj[queryArr[0]] = queryArr[1];
  });
  let gr = obj["grade"];
  let cl = obj["class"];
  let ki = obj["kind"];
  request.onreadystatechange = function() {
    if (request.readyState == 4 && request.status == 200) {
      let text = JSON.parse(request.responseText);
      search(text,gr,cl,ki);
    }
  };
  request.send(null);
});

function search(data,gr,cl,ki){
  let result = document.getElementById("result");
  for(let i of data["tags"]){
    console.log(gr);
    if(i["grade"]==gr || i["class"]==cl || i["kind"]==ki){
      result.insertAdjacentHTML("beforeend",data["text"][i["pointer"]]+"<br><hr><br>");
    }
  }
}
