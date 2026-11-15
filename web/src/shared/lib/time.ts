// Seconds as m:ss; negative input clamps to 0:00.
export const fmtClock = (sec: number) => {
  const t = Math.max(0, Math.floor(sec))
  return `${Math.floor(t / 60)}:${String(t % 60).padStart(2, '0')}`
}
