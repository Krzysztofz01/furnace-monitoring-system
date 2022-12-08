export function isDateCurrentDay(date: Date): boolean {
    const dateObj = new Date(date);
    const today = new Date(Date.now());

    return dateObj.getDate() === today.getDate() &&
        dateObj.getMonth() === today.getMonth() &&
        dateObj.getFullYear() === today.getFullYear();
}