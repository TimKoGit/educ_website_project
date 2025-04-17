import { combineReducers } from 'redux';
import contestsReducer from './slices/contestsSlice';
import authReducer from './slices/authSlice';
import submissionsReducer from './slices/submissionsSlice';
import tasksReducer from './slices/tasksSlice';


const rootReducer = combineReducers({
  auth: authReducer,
  contests: contestsReducer,
  submissions: submissionsReducer,
  tasks: tasksReducer
});


export default rootReducer;
