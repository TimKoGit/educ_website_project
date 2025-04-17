import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import axios from 'axios';


interface Contest {
  id: string;
  name: string;
  time: string;
  duration: string;
  groupid: string;
}


interface ContestsState {
  contests: Contest[] | null;
  contestsById: Record<string, Contest>;
  fetchContestsByGroupIdLoading: boolean;
  fetchContestsByGroupIdError: string | null;
  addNewContestLoading: boolean;
  addNewContestError: string | null;
  fetchContestLoading: boolean;
  fetchContestError: string | null;
  deleteContestLoading: boolean;
  deleteContestError: string | null;
}


const initialState: ContestsState = {
  contests: [],
  contestsById: {},
  fetchContestsByGroupIdLoading: false,
  fetchContestsByGroupIdError: null,
  addNewContestLoading: false,
  addNewContestError: null,
  fetchContestLoading: false,
  fetchContestError: null,
  deleteContestLoading: false,
  deleteContestError: null,
};


export const fetchContestsByGroupId = createAsyncThunk(
  'contests/fetchContestsByGroupId',
  async ({ groupid }: { groupid: string }, { rejectWithValue }) => {
    try {
      const response = await axios.get(`http://130.193.44.85:5000/contests_by_groupid?groupid=${groupid}`, {withCredentials: true});
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error || 'Не удалось получить контест');
    } 
  }
);


export const addNewContest = createAsyncThunk(
  'contests/addNewContest',
  async (newContest: { name: string; groupid: string, time: string, duration: string }, { rejectWithValue }) => {
    try {
      const response = await axios.post('http://130.193.44.85:5000/contests', newContest, {withCredentials: true});
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Не удалось добавить контест');
    }
  }
);


export const fetchContestById = createAsyncThunk(
  'contests/fetchContestById',
  async (contestid: string, { rejectWithValue }) => {
    try {
      const response = await axios.get(`http://130.193.44.85:5000/contests/${contestid}`, {withCredentials: true});
      return response.data as Contest;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Не удалось получить контест');
    }
  }
);


export const deleteContest = createAsyncThunk(
  'contests/deleteContest',
  async (contestid: string, { rejectWithValue }) => {
    try {
      await axios.delete(`http://130.193.44.85:5000/contests/${contestid}`, {withCredentials: true});
      return contestid;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Не удалось удалить контест');
    }
  }
);


const contestsSlice = createSlice({
  name: 'contests',
  initialState,
  reducers: {
    clearContests(state) {
      state.contests = [];
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchContestsByGroupId.pending, (state) => {
        state.fetchContestsByGroupIdLoading = true;
        state.fetchContestsByGroupIdError = null;
      })
      .addCase(fetchContestsByGroupId.fulfilled, (state, action) => {
        state.fetchContestsByGroupIdLoading = false;
        state.contests = action.payload;
      })
      .addCase(fetchContestsByGroupId.rejected, (state, action) => {
        state.fetchContestsByGroupIdLoading = false;
        state.fetchContestsByGroupIdError = action.payload as string;
      })
      .addCase(addNewContest.pending, (state) => {
        state.addNewContestLoading = true;
        state.addNewContestError = null;
      })
      .addCase(addNewContest.fulfilled, (state, action) => {
        state.addNewContestLoading = false;
      })
      .addCase(addNewContest.rejected, (state, action) => {
        state.addNewContestLoading = false;
        state.addNewContestError = action.payload as string;
      })
      .addCase(fetchContestById.pending, (state) => {
        state.fetchContestLoading = true;
        state.fetchContestError = null;
      })
      .addCase(fetchContestById.fulfilled, (state, action: PayloadAction<Contest>) => {
        state.fetchContestLoading = false;
        state.contestsById[action.payload.id] = action.payload;
      })
      .addCase(fetchContestById.rejected, (state, action) => {
        state.fetchContestLoading = false;
        state.fetchContestError = action.payload as string;
      })
      .addCase(deleteContest.pending, (state) => {
        state.deleteContestLoading = true;
        state.deleteContestError = null;
      })
      .addCase(deleteContest.fulfilled, (state, action: PayloadAction<string>) => {
        state.deleteContestLoading = false;
        state.contests = state.contests.filter(contest => contest.id !== action.payload);
      })
      .addCase(deleteContest.rejected, (state, action) => {
        state.deleteContestLoading = false;
        state.deleteContestError = action.payload as string;
      });
  },
});


export const { clearContests } = contestsSlice.actions;
export default contestsSlice.reducer;
