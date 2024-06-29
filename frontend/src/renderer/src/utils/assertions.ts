export function ensureNonEmpty(value : string) {
    if (value ===undefined) {
        throw new Error('Expected string to be defined')
    }
    if (value === null) {
        throw new Error('Expected string to be not null')
    }
    if (value === '') {
        throw new Error('Expected string to be not empty')
    }
}