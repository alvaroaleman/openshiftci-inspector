import {Job, JobsApi} from "../api-client";
import INotificationService, {NotificationType} from "../notification/service/INotificationService";

export default class JobsGetRelatedService {

    constructor(
        private api: JobsApi,
        private notificationService: INotificationService
    ) {}

    public async getJob(id: string): Promise<Array<Job>> {
        try {
            let response = await this.api.getRelatedJobs(id)
            return response.data.jobs
        } catch (e) {
            this.notificationService.notifyUser({
                type: NotificationType.ERROR,
                message: "Failed to fetch related jobs for " + id + "."
            })
            console.log(e)
            throw e
        }
    }
}