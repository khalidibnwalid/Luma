
enum METHOD {
    POST = 'POST',
    PUT = 'PUT',
    PATCH = 'PATCH',
    DELETE = 'DELETE'
}

const JWT_LOCAL_STORAGE_KEY = 'token'

export default function http(url: string) {
    const isFormData = (body: unknown) => body instanceof FormData
    const get = async <T>() => (await fetchWrapper(url)).json() as T
    const post = async <T>(body?: unknown) => (await apiFactory(METHOD.POST)(body)).json() as T
    const patch = async <T>(body?: unknown) => (await apiFactory(METHOD.PATCH)(body)).json() as T
    const put = async <T>(body?: unknown) => (await apiFactory(METHOD.PUT)(body)).json() as T
    const del = async <T>(body?: unknown) => (await apiFactory(METHOD.DELETE)(body)).json() as T

    function apiFactory(method: METHOD) {
        return (body?: unknown) =>
        (fetchWrapper(url, {
            method,
            headers: {
                'Content-Type': isFormData(body) ? 'multipart/form-data' : 'application/json'
            },
            body: isFormData(body) ? body : JSON.stringify(body)
        }))
    }

    return {
        get,
        post,
        patch,
        put,
        delete: del
    }
}

async function fetchWrapper(url: string, requestInit?: RequestInit) {
    try {
        const res = await fetch(url, {
            credentials: 'include',
            headers: {
                Authorization: `Bearer ${localStorage.getItem(JWT_LOCAL_STORAGE_KEY)}`,
            },
            ...requestInit
        })
        const clone = res.clone() // Clone the response to check for errors, since res.json() consumes it
        const jsonRes = await res.json();

        if (!res.ok)
            throw new Error(jsonRes.message)

        return clone
    } catch (e) {
        throw new Error((e as Error).message)
    }
}