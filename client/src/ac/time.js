export function humanDuration(seconds) {
  const h = Math.floor(seconds / 3600);
  const m = Math.floor(seconds % 3600 / 60);
  const s = Math.floor(seconds % 3600 % 60);

  if (h + m + s === 0) {
    return '-';
  }

  const hs = h > 0 ? h + 'h ' : '';
  const ms = h > 0 || m > 0 ? padZero(m) + 'm ' : '';
  const ss = h > 0 || m > 0 || s > 0 ? padZero(s) + 's' : '';

  return hs + ms + ss;
}

function padZero(num) {
  return num < 10 ? '0' + num.toString() : num;
}
