import {Job, JobsApi} from "../api-client";
import {AbstractService} from "../service/service";
import INotificationService, {NotificationType} from "../notification/service/INotificationService";

export default class JobsListService extends AbstractService {
    private loaded: boolean
    private refreshing: boolean
    private jobs: Array<Job> = new Array<Job>()

    constructor(
        private api: JobsApi,
        private notificationService: INotificationService
    ) {
        super()
        this.loaded = false
        this.refreshing = false
        this.jobs = new Array<Job>()
    }

    public getJobs(): Array<Job> {
        return this.jobs
    }

    public isLoaded(): boolean {
        return this.loaded
    }

    public isRefreshing(): boolean{
        return this.refreshing
    }

    public async refresh() {
        this.refreshing = true
        this.notify()
        try {
            let response = await this.api.listJobs()
            if (response.data.jobs) {
                this.jobs = response.data.jobs
            }
            this.loaded = true
        } catch (e) {
            this.notificationService.notifyUser({
                type: NotificationType.ERROR,
                message: "Failed to update jobs list."
            })
            console.log(e)
        } finally {
            this.refreshing = false
            this.notify()
        }
    }
}