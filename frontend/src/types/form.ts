export type OldInputValue = string | string[]

export type OldInput = Record<string, OldInputValue>

export type ValidationErrors = Record<string, string>

export function getOldInputString(oldInput: OldInput, key: string) {
  const value = oldInput[key]
  return typeof value === 'string' ? value : ''
}

export function getOldInputStringList(oldInput: OldInput, key: string) {
  const value = oldInput[key]
  if (Array.isArray(value)) {
    return value
  }
  return typeof value === 'string' && value !== '' ? [value] : []
}
