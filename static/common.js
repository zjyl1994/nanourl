const $ = document.querySelector.bind(document);
const $$ = document.querySelectorAll.bind(document);

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

function get_page_param(param_name, default_value) {
    let params = (new URL(document.location)).searchParams;
    let data = params.get(param_name);
    return data ? data : default_value;
}

function get_pagination(page, page_size, total_rows) {
    const total_page = Math.ceil(total_rows / page_size);
    const prev_flag = page > 1;
    const next_flag = page < total_page;
    var page_arr = [];
    if (prev_flag)
        page_arr.push(page - 1);
    page_arr.push(page);
    if (next_flag)
        page_arr.push(page + 1);
    return {
        prev: prev_flag,
        next: next_flag,
        nums: page_arr,
        curr: page,
    }
}