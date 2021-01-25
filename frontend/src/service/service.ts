
export interface IServiceHandler {
    update():void;
}

export interface IService {
    register(handler: IServiceHandler): void;
    deregister(handler: IServiceHandler): void;
}

export abstract class AbstractService implements IService{
    private handlers: Set<IServiceHandler> = new Set();

    deregister(handler: IServiceHandler): void {
        this.handlers.delete(handler)
    }

    register(handler: IServiceHandler): void {
        this.handlers.add(handler)
    }

    notify() {
        this.handlers.forEach(handler => handler.update())
    }
}