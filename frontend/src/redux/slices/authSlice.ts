import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import axios from 'axios';


interface User {
  id: string,
  role: string,
  surname: string,
  firstname: string,
}


interface Group {
  id: string;
  name: string;
}


interface AuthState {
  user: User | null;
  groups: Group[] | null;
  loginLoading: boolean;
  loginError: string | null;
  registerLoading: boolean;
  registerError: string | null;
  fetchGroupsByUserIdLoading: boolean;
  fetchGroupsByUserIdError: string | null;
  addNewGroupLoading: boolean;
  addNewGroupError: string | null;
  joinGroupByCodeLoading: boolean;
  joinGroupByCodeError: string | null;
  deleteGroupLoading: boolean;
  deleteGroupError: string | null;
}


const storedUser = localStorage.getItem('user');
const initialState: AuthState = {
  user: storedUser ? JSON.parse(storedUser) : null,
  groups: null,
  loginLoading: false,
  loginError: null,
  registerLoading: false,
  registerError: null,
  fetchGroupsByUserIdLoading: false,
  fetchGroupsByUserIdError: null,
  addNewGroupLoading: false,
  addNewGroupError: null,
  joinGroupByCodeLoading: false,
  joinGroupByCodeError: null,
  deleteGroupLoading: false,
  deleteGroupError: null,
};


export const login = createAsyncThunk(
  'auth/login',
  async (credentials: { username: string; password: string }, { rejectWithValue }) => {
    try {
      const response = await axios.get(`http://localhost:5000/login`, {
        params: {
          username: credentials.username,
          password: credentials.password,
        },
        withCredentials: true
      });

      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error || 'Не удалось войти');
    }
  }
);


export const register = createAsyncThunk(
  'auth/register',
  async (credentials: { username: string; password: string; firstname: string; surname: string; role: string }, { rejectWithValue }) => {
    try {
      debugger;
      const response = await axios.post(
        'http://localhost:5000/register',
        null,
        {
          params: {
            username: credentials.username,
            password: credentials.password,
            firstname: credentials.firstname,
            surname: credentials.surname,
            role: credentials.role,
          },
          withCredentials: true,
        }
      );

      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error || 'Не удалось зарегистрироваться');
    }
  }
);


export const fetchGroupsByUserId = createAsyncThunk(
  'auth/fetchGroupsByUserId',
  async (_, { rejectWithValue }) => {
    try {
      const userGroupsResponse = await axios.get('http://localhost:5000/groups', { withCredentials: true });
      return userGroupsResponse.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error || 'Не удалось получить группы пользователя');
    }
  }
);


export const addNewGroup = createAsyncThunk(
  'groups/addNewGroup',
  async (newGroup: { name: string; code: string }, {rejectWithValue}) => {
    try {
      const groupResponse = await axios.post('http://localhost:5000/groups', {
        name: newGroup.name,
        code: newGroup.code,
      }, { withCredentials: true });
      return groupResponse.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error || 'Не удалось добавить группу');
    }
  }
);


export const joinGroupByCode = createAsyncThunk(
  'auth/joinGroupByCode',
  async (groupCode: string, { rejectWithValue }) => {
    try {
      const joinGroupResponse = await axios.post('http://localhost:5000/join_group', {
        code: groupCode
      }, { withCredentials: true });
      return joinGroupResponse.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.error || 'Не удалось присоединиться к группе');
    }
  }
);


export const deleteGroup = createAsyncThunk(
  'groups/deleteGroup',
  async (groupid: string, { rejectWithValue }) => {
    try {
      await axios.delete(`http://localhost:5000/groups/${groupid}`, {withCredentials: true});
      return groupid;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Не удалось удалить группу');
    }
  }
);


const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    logout(state) {
      state.user = null;
      state.groups = null;
      localStorage.removeItem('user');
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(login.pending, (state) => {
        state.loginLoading = true;
        state.loginError = null;
      })
      .addCase(login.fulfilled, (state, action) => {
        state.loginLoading = false;
        state.user = action.payload;
        localStorage.setItem('user', JSON.stringify(action.payload));
      })
      .addCase(login.rejected, (state, action) => {
        state.loginLoading = false;
        state.loginError = action.payload as string;
      })
      .addCase(register.pending, (state) => {
        state.registerLoading = true;
        state.registerError = null;
      })
      .addCase(register.fulfilled, (state, action) => {
        state.registerLoading = false;
        state.user = action.payload;
        localStorage.setItem('user', JSON.stringify(action.payload));
      })
      .addCase(register.rejected, (state, action) => {
        state.registerLoading = false;
        state.registerError = action.payload as string;
      })
      .addCase(fetchGroupsByUserId.pending, (state) => {
        state.fetchGroupsByUserIdLoading = true;
        state.fetchGroupsByUserIdError = null;
      })
      .addCase(fetchGroupsByUserId.fulfilled, (state, action) => {
        state.fetchGroupsByUserIdLoading = false;
        state.groups = action.payload;
      })
      .addCase(fetchGroupsByUserId.rejected, (state, action) => {
        state.fetchGroupsByUserIdLoading = false;
        state.fetchGroupsByUserIdError = action.payload as string;
      })
      .addCase(addNewGroup.pending, (state) => {
        state.addNewGroupLoading = true;
        state.addNewGroupError = null;
      })
      .addCase(addNewGroup.fulfilled, (state, action) => {
        state.addNewGroupLoading = false;
      })
      .addCase(addNewGroup.rejected, (state, action) => {
        state.addNewGroupLoading = false;
        state.addNewGroupError = action.payload as string;
      })
      .addCase(joinGroupByCode.pending, (state) => {
        state.joinGroupByCodeLoading = true;
        state.joinGroupByCodeError = null;
      })
      .addCase(joinGroupByCode.fulfilled, (state, action) => {
        state.joinGroupByCodeLoading = false;
        if (state.groups) {
          state.groups.push(action.payload);
        } else {
          state.groups = [action.payload];
        }
      })
      .addCase(joinGroupByCode.rejected, (state, action) => {
        state.joinGroupByCodeLoading = false;
        state.joinGroupByCodeError = action.payload as string;
      })
      .addCase(deleteGroup.pending, (state) => {
        state.deleteGroupLoading = true;
        state.deleteGroupError = null;
      })
      .addCase(deleteGroup.fulfilled, (state, action: PayloadAction<string>) => {
        state.deleteGroupLoading = false;
        state.groups = state.groups.filter(group => group.id !== action.payload);
      })
      .addCase(deleteGroup.rejected, (state, action) => {
        state.deleteGroupLoading = false;
        state.deleteGroupError = action.payload as string;
      });
  },
});


export const { logout } = authSlice.actions;
export default authSlice.reducer;
