export enum DockerContainerInfoState {
  created = 'created',
  restarting = 'restarting',
  removing = 'removing',
  running = 'running',
  paused = 'paused',
  exited = 'exited',
  dead = 'dead',
  error = 'error', // 自定义状态与原生docker无关
}
