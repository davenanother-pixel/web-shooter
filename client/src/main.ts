const canvas = document.getElementById('c') as HTMLCanvasElement;
const ctx    = canvas.getContext('2d')!;
const sock   = new WebSocket('ws://localhost:8080/ws');

interface State { x: number; y: number; r: number; rgb: [number,number,number] }
let state: State = { x: 400, y: 300, r: 20, rgb: [0,255,255] };

sock.onmessage = e => { state = JSON.parse(e.data); draw(); };
function draw(){
  ctx.fillStyle = '#111';
  ctx.fillRect(0,0,canvas.width,canvas.height);
  const [r,g,b] = state.rgb;
  ctx.beginPath();
  ctx.arc(state.x, state.y, state.r, 0, Math.PI*2);
  ctx.fillStyle = `rgb(${r},${g},${b})`;
  ctx.fill();
}

window.addEventListener('keydown', (ev)=>{
  const code = ev.key === 'ArrowUp' ? 'U' : ev.key === 'ArrowDown' ? 'D' :
               ev.key === 'ArrowLeft' ? 'L' : ev.key === 'ArrowRight' ? 'R' : '';
  if(code) sock.send(JSON.stringify({cmd: code}));
});
