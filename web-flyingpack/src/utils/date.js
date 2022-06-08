
//currentDate return string date in format yyyy-mm-dd
export const currentDate = () => {
    const now = new Date()
    const date = String(now.getDate()).padStart(2, '0')
    const month = String(now.getMonth()+1).padStart(2, '0')
    const year = String(now.getFullYear())

    return `${year}-${month}-${date}`
}
