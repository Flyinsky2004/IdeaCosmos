import {defineStore} from "pinia";

const useProjectStore = defineStore('project', {
    state: () => ({
        project: null,
    }),

    actions: {
        setProject(project){
            this.project = project
        }
    },

    getters: {
        isSelectedProject: (state) => state.project !== null
    }
})
export { useProjectStore }