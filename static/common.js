function format_ts(timestamp) {
    const pad2 = (num) => num < 10 ? '0' + num : num;
    let date = new Date(timestamp * 1000);
    let year = date.getFullYear();
    let month = pad2(date.getMonth() + 1);
    let day = pad2(date.getDate());
    let hour = pad2(date.getHours());
    let minute = pad2(date.getMinutes());
    let second = pad2(date.getSeconds());
    return year + '-' + month + '-' + day + ' ' + hour + ':' + minute + ':' + second;
}