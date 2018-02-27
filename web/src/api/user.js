import request from '@/utils/request'

export function changePasswd(data) {
  console.log(data)
  return request({
    url: '/user/changePasswd',
    method: 'post',
    data: data
  })
}

export function userList(data) {
  console.log(data)
  return request({
    url: '/user/userList',
    method: 'post',
    data: data
  })
}

export function edit(data) {
  console.log(data)
  return request({
    url: '/user/edit',
    method: 'post',
    data: data
  })
}
