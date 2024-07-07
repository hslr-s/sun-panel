export function verifyUrlHttp(url: string): boolean {
  if (url === '')
    return true

  return url.startsWith('https://') || url.startsWith('http://')
}
