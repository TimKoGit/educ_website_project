import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import axios from 'axios';


interface Submission {
  id: string;
  studentName: string;
  taskid: string;
  taskName: string;
  groupid: string;
  contestid: string;
  status: 'unchecked' | 'accepted' | 'declined';
}


interface SubmissionPhoto {
  id: string;
  url: string;
  studentName: string;
  taskName: string;
}


interface SubmissionsState {
  uncheckedSubmissions: Submission[] | null;
  specificSubmission: SubmissionPhoto | null;
  fetchUncheckedSubmissionsByGroupIdLoading: boolean;
  fetchUncheckedSubmissionsByGroupIdError: string | null;
  fetchUncheckedSubmissionsByContestIdLoading: boolean;
  fetchUncheckedSubmissionsByContestIdError: string | null;
  fetchSubmissionByIdLoading: boolean;
  fetchSubmissionByIdError: string | null;
  addSubmissionLoading: boolean;
  addSubmissionError: string | null;
  acceptSubmissionLoading: boolean;
  acceptSubmissionError: string | null;
  declineSubmissionLoading: boolean;
  declineSubmissionError: string | null;
}


const initialState: SubmissionsState = {
  uncheckedSubmissions: [],
  specificSubmission: null,
  fetchUncheckedSubmissionsByGroupIdLoading: false,
  fetchUncheckedSubmissionsByGroupIdError: null,
  fetchUncheckedSubmissionsByContestIdLoading: false,
  fetchUncheckedSubmissionsByContestIdError: null,
  fetchSubmissionByIdLoading: false,
  fetchSubmissionByIdError: null,
  addSubmissionLoading: false,
  addSubmissionError: null,
  acceptSubmissionLoading: false,
  acceptSubmissionError: null,
  declineSubmissionLoading: false,
  declineSubmissionError: null,
};


export const fetchUncheckedSubmissionsByGroupId = createAsyncThunk(
  'submissions/fetchUncheckedSubmissionsByGroupId',
  async ({ groupid }: { groupid: string }, { rejectWithValue }) => {
    try {
      const response = await axios.get(`http://130.193.44.85:5000/submissions/unchecked_by_groupid/${groupid}`, {withCredentials: true});
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Не удалось получить посылки');
    }
  }
);


export const fetchUncheckedSubmissionsByContestId = createAsyncThunk(
  'submissions/fetchUncheckedSubmissionsByContestId',
  async ({ contestid }: { contestid: string }, { rejectWithValue }) => {
    try {
      const response = await axios.get(`http://130.193.44.85:5000/submissions/unchecked_by_contestid/${contestid}`, {withCredentials: true});
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Не удалось получить посылки');
    }
  }
);


export const fetchSubmissionById = createAsyncThunk(
  'submissions/fetchSubmissionById',
  async ({submissionid}: {submissionid: string}, { rejectWithValue }) => {
    try {
      const response_picture = await axios.get(`http://130.193.44.85:5000/submissions_picture/${submissionid}`,
      {withCredentials: true, responseType: 'blob'});
      const imageUrl = URL.createObjectURL(response_picture.data);
      const response_details = await axios.get(`http://130.193.44.85:5000/submissions_details/${submissionid}`, {
        withCredentials: true,
      });
      const { studentName, taskName } = response_details.data;
        
      return { id: submissionid, url: imageUrl, studentName, taskName } as SubmissionPhoto;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Не удалось получить посылку');
    }
  }
);


export const addSubmission = createAsyncThunk(
  'submissions/addSubmission',
  async (newSubmission: { studentName: string; taskid: string; groupid: string }, { rejectWithValue }) => {
    try {
      const response = await axios.post('http://130.193.44.85:5000/submissions', newSubmission);
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.message || 'Не удалось добавить посылку');
    }
  }
);


export const acceptSubmission = createAsyncThunk(
  'submissions/acceptSubmission',
  async (submissionId: string, { rejectWithValue }) => {
    try {
      const response = await axios.patch(`http://130.193.44.85:5000/submissions/accept/${submissionId}`, {}, {withCredentials: true});
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Не удалось принять посылку');
    }
  }
);


export const declineSubmission = createAsyncThunk(
  'submissions/declineSubmission',
  async (submissionId: string, { rejectWithValue }) => {
    try {
      const response = await axios.patch(`http://130.193.44.85:5000/submissions/decline/${submissionId}`, {}, {withCredentials: true});
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Не удалось отклонить посылку');
    }
  }
);


const submissionsSlice = createSlice({
  name: 'submissions',
  initialState,
  reducers: {
    clearSubmissions(state) {
      state.uncheckedSubmissions = [];
      state.fetchUncheckedSubmissionsByGroupIdError = null;
      state.fetchUncheckedSubmissionsByGroupIdLoading = false;
    },
    clearSpecificSubmission(state) {
      state.specificSubmission = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchUncheckedSubmissionsByGroupId.pending, (state) => {
        state.fetchUncheckedSubmissionsByGroupIdLoading = true;
        state.fetchUncheckedSubmissionsByGroupIdError = null;
      })
      .addCase(fetchUncheckedSubmissionsByGroupId.fulfilled, (state, action: PayloadAction<Submission[]>) => {
        state.fetchUncheckedSubmissionsByGroupIdLoading = false;
        state.uncheckedSubmissions = action.payload;
      })
      .addCase(fetchUncheckedSubmissionsByGroupId.rejected, (state, action) => {
        state.fetchUncheckedSubmissionsByGroupIdLoading = false;
        state.fetchUncheckedSubmissionsByGroupIdError = action.payload as string;
      })
      .addCase(fetchUncheckedSubmissionsByContestId.pending, (state) => {
        state.fetchUncheckedSubmissionsByContestIdLoading = true;
        state.fetchUncheckedSubmissionsByContestIdError = null;
      })
      .addCase(fetchUncheckedSubmissionsByContestId.fulfilled, (state, action: PayloadAction<Submission[]>) => {
        state.fetchUncheckedSubmissionsByContestIdLoading = false;
        state.uncheckedSubmissions = action.payload;
      })
      .addCase(fetchUncheckedSubmissionsByContestId.rejected, (state, action) => {
        state.fetchUncheckedSubmissionsByContestIdLoading = false;
        state.fetchUncheckedSubmissionsByContestIdError = action.payload as string;
      })
      .addCase(fetchSubmissionById.pending, (state) => {
        state.fetchSubmissionByIdLoading = true;
        state.fetchSubmissionByIdError = null;
        state.specificSubmission = null;
      })
      .addCase(fetchSubmissionById.fulfilled, (state, action: PayloadAction<SubmissionPhoto>) => {
        state.fetchSubmissionByIdLoading = false;
        state.specificSubmission = action.payload;
      })
      .addCase(fetchSubmissionById.rejected, (state, action) => {
        state.fetchSubmissionByIdLoading = false;
        state.fetchSubmissionByIdError = action.payload as string;
      })
      .addCase(acceptSubmission.pending, (state) => {
        state.acceptSubmissionLoading = true;
        state.acceptSubmissionError = null;
      })
      .addCase(acceptSubmission.fulfilled, (state, action: PayloadAction<Submission>) => {
        state.acceptSubmissionLoading = false;
      })
      .addCase(acceptSubmission.rejected, (state, action) => {
        state.acceptSubmissionLoading = false;
        state.acceptSubmissionError = action.payload as string;
      })
      .addCase(declineSubmission.pending, (state) => {
        state.declineSubmissionLoading = true;
        state.declineSubmissionError = null;
      })
      .addCase(declineSubmission.fulfilled, (state, action: PayloadAction<Submission>) => {
        state.declineSubmissionLoading = false;
      })
      .addCase(declineSubmission.rejected, (state, action) => {
        state.declineSubmissionLoading = false;
        state.declineSubmissionError = action.payload as string;
      })
      .addCase(addSubmission.pending, (state) => {
        state.addSubmissionLoading = true;
        state.addSubmissionError = null;
      })
      .addCase(addSubmission.fulfilled, (state, action: PayloadAction<Submission>) => {
        state.addSubmissionLoading = false;
      })
      .addCase(addSubmission.rejected, (state, action) => {
        state.addSubmissionLoading = false;
        state.addSubmissionError = action.payload as string;
      });
  },
});


export const { clearSubmissions, clearSpecificSubmission } = submissionsSlice.actions;
export default submissionsSlice.reducer;
