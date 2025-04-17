import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import axios from 'axios';


interface Task {
  id: string;
  name: string;
  url: string;
  status: string;
}


interface TasksState {
  tasks: Task[];
  fetchTasksLoading: boolean;
  fetchTasksError: string | null;
  deleteTaskLoading: boolean;
  deleteTaskError: string | null;
  addTaskLoading: boolean;
  addTaskError: string | null;
}


const initialState: TasksState = {
  tasks: [],
  fetchTasksLoading: false,
  fetchTasksError: null,
  deleteTaskLoading: false,
  deleteTaskError: null,
  addTaskLoading: false,
  addTaskError: null,
};


export const fetchTasksByContestId = createAsyncThunk(
  'tasks/fetchTasksByContestId',
  async (contestid: string, { rejectWithValue }) => {
    try {
      const response = await axios.get(`http://localhost:5000/tasks/${contestid}`, {withCredentials: true});
      debugger;
      return response.data as Task[];
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Не удалось получить задачи');
    }
  }
);


export const deleteTask = createAsyncThunk(
  'tasks/deleteTask',
  async (taskid: string, { rejectWithValue }) => {
    try {
      await axios.delete(`http://localhost:5000/tasks/${taskid}`, {withCredentials: true});
      return taskid;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Не удалось удалить задачу');
    }
  }
);


export const addTask = createAsyncThunk(
  'tasks/addTask',
  async (newTask: {name: string, url: string, contestid: string}, { rejectWithValue }) => {
    try {
      const response = await axios.post(`http://localhost:5000/tasks`, newTask, {withCredentials: true});
      return response.data as Task;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Не удалось добваить задачу');
    }
  }
);


const tasksSlice = createSlice({
  name: 'tasks',
  initialState,
  reducers: {
    clearTasks(state) {
      state.tasks = [];
      state.fetchTasksError = null;
      state.fetchTasksLoading = false;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchTasksByContestId.pending, (state) => {
        state.fetchTasksLoading = true;
        state.fetchTasksError = null;
      })
      .addCase(fetchTasksByContestId.fulfilled, (state, action: PayloadAction<Task[]>) => {
        state.fetchTasksLoading = false;
        state.tasks = action.payload;
      })
      .addCase(fetchTasksByContestId.rejected, (state, action) => {
        state.fetchTasksLoading = false;
        state.fetchTasksError = action.payload as string;
      })
      .addCase(deleteTask.pending, (state) => {
        state.deleteTaskLoading = true;
        state.deleteTaskError = null;
      })
      .addCase(deleteTask.fulfilled, (state, action: PayloadAction<string>) => {
        state.deleteTaskLoading = false;
        state.tasks = state.tasks.filter(task => task.id !== action.payload);
      })
      .addCase(deleteTask.rejected, (state, action) => {
        state.deleteTaskLoading = false;
        state.deleteTaskError = action.payload as string;
      })
      .addCase(addTask.pending, (state) => {
        state.addTaskLoading = true;
        state.addTaskError = null;
      })
      .addCase(addTask.fulfilled, (state, action: PayloadAction<Task>) => {
        state.addTaskLoading = false;
      })
      .addCase(addTask.rejected, (state, action) => {
        state.addTaskLoading = false;
        state.addTaskError = action.payload as string;
      });;
  },
});


export const { clearTasks } = tasksSlice.actions;
export default tasksSlice.reducer;
