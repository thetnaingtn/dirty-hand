import { createChannel, createClientFactory, FetchTransport } from 'nice-grpc-web'
import { ProductServiceDefinition } from './types/proto/api/v1/product'
import { UserServiceDefinition } from './types/proto/api/v1/user'

const channel = createChannel(
    window.location.origin,
    FetchTransport({
        credentials:"include"
    })
)

const clientFactory = createClientFactory()

export const productClient = clientFactory.create(ProductServiceDefinition, channel)
export const userClient = clientFactory.create(UserServiceDefinition, channel)