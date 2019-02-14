var message={
  name:"channel add",
  data:{
    name:"Hardware Support"
  }
}
var submsg={
  name:"channel subscribe"
}

var ws=new WebSocket("ws://localhost:4001")

ws.onopen = ()=>{
  ws.send(JSON.stringify(submsg));
  ws.send(JSON.stringify(message));
  
}
ws.onmessage=(data)=>{
  console.log(data)
}