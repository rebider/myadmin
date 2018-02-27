import request from '@/utils/request'

export function roleList(data) {
  console.log(data)
  return request({
    url: '/role/roleList',
    method: 'post',
    data: data
  })
}
