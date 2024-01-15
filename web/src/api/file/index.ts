import request from '@/api/axiosInstance';

export function getFiles() {
    return request({
        url: '/file/all',
        method: 'get',
    });
}
