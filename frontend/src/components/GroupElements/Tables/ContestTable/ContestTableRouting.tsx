import * as React from 'react';
import { useParams } from 'react-router-dom';
import { useEffect } from 'react';
import { useAppSelector, useAppDispatch } from '../../../../redux/hooks';
import { fetchContestsByGroupId } from '../../../../redux/slices/contestsSlice';
import { fetchTasksByContestId, clearTasks } from '../../../../redux/slices/tasksSlice';
import ContestTable from './ContestTable';
import TaskTable from './TaskTable';


const ContestTableRouting: React.FC = () => {
  const { groupid, contestid } = useParams<{ groupid: string; contestid?: string }>();
  const dispatch = useAppDispatch();
  const { contests, fetchContestsByGroupIdLoading, fetchContestsByGroupIdError } = useAppSelector((state) => state.contests);
  const { tasks, fetchTasksLoading, fetchTasksError } = useAppSelector((state) => state.tasks);

  const user = useAppSelector((state) => state.auth.user);

  useEffect(() => {
    if (groupid && user) {
      dispatch(fetchContestsByGroupId({ groupid }));
    }
  }, [dispatch, groupid, user]);

  useEffect(() => {
    if (contestid) {
      dispatch(fetchTasksByContestId(contestid));
    } else {
      dispatch(clearTasks());
    }
  }, [dispatch, contestid]);

  if (!user || (!groupid && !contestid)) {
    return <div />;
  } 

  if (fetchContestsByGroupIdLoading || fetchTasksLoading) {
    return <div>Загрузка...</div>;
  }

  if (fetchContestsByGroupIdError) {
    return <div>Ошибка: {fetchContestsByGroupIdError}</div>;
  }

  if (fetchTasksError) {
    return <div>Ошибка: {fetchTasksError}</div>;
  }

  if (contestid) {
    return <TaskTable tasks={tasks} />;
  }

  return <ContestTable contest_infos={contests} />;
};


export default ContestTableRouting;
