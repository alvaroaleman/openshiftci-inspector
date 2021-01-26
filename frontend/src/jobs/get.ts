import {JobsApi, JobWithAssetURL} from "../api-client";
import INotificationService, {NotificationType} from "../notification/service/INotificationService";

export default class JobsGetService {

    constructor(
        private api: JobsApi,
        private notificationService: INotificationService
    ) {}

    public async getJob(id: string): Promise<JobWithAssetURL> {
        try {
            let response = await this.api.getJob(id)
            return response.data.job
        } catch (e) {
            this.notificationService.notifyUser({
                type: NotificationType.ERROR,
                message: "Failed to fetch job " + id + "."
            })
            console.log(e)
            throw e
        }
    }
}