let radios = document.querySelectorAll('.radio');
let labels = document.querySelectorAll('.label');
let ball = document.getElementById('gball');
let prevRadio, prevLabel;
radios.forEach((radio, index) => {
  radio.addEventListener('click', function(e) {
    if (prevRadio) prevRadio.classList.toggle('active');
    if (prevLabel) prevLabel.classList.toggle('active');
    radio.classList.toggle('active');
    prevRadio = radio;
    labels[index].classList.toggle('active');
    prevLabel = labels[index];
    ball.className = `ball pos${index}`;
  });
});

$(function(){
  $('#toclass').on("click",function(){
    $('#grade').css({
      "display":"none"
    });
    $('#class').css({
      "display":"flex"
    });
    ball = document.getElementById('cball');
  });
  $('#tokind').on("click",function(){
    $('#class').css({
      "display":"none"
    });
    $('#kind').css({
      "display":"flex"
    });
    ball = document.getElementById('kball');
  });
  $('#toinput').on("click",function(){
    $('#kind').css({
      "display":"none"
    });
    $('#text').css({
      "display":"flex"
    });
    ball = document.getElementById('kball');
  });
});