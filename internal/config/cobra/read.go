package cobra

// func createNotificationSettings(cmd *cobra.Command) (*types.NotificationSettings, *errors.Error) {
// 	var err *errors.Error

// 	var notificationUrls []string
// 	if notificationUrls, err = StringSlice(keys.NotificationUrls, cmd.Flags()); err != nil {
// 		return nil, err
// 	}

// 	var notifyJobStarted bool
// 	if notifyJobStarted, err = Bool(keys.NotifyJobStarted, cmd.Flags()); err != nil {
// 		return nil, err
// 	}

// 	var notifyJobCompleted bool
// 	if notifyJobCompleted, err = Bool(keys.NotifyJobCompleted, cmd.Flags()); err != nil {
// 		return nil, err
// 	}

// 	var notifyJobError bool
// 	if notifyJobError, err = Bool(keys.NotifyJobError, cmd.Flags()); err != nil {
// 		return nil, err
// 	}

// 	var notifyTaskStarted bool
// 	if notifyTaskStarted, err = Bool(keys.NotifyTaskStarted, cmd.Flags()); err != nil {
// 		return nil, err
// 	}

// 	var notifyTaskCompleted bool
// 	if notifyTaskCompleted, err = Bool(keys.NotifyTaskCompleted, cmd.Flags()); err != nil {
// 		return nil, err
// 	}

// 	var notifyTaskError bool
// 	if notifyTaskError, err = Bool(keys.NotifyTaskError, cmd.Flags()); err != nil {
// 		return nil, err
// 	}

// 	return &types.NotificationSettings{
// 		NotificationUrls: notificationUrls,

// 		NotifyJobStarted:    notifyJobStarted,
// 		NotifyJobCompleted:  notifyJobCompleted,
// 		NotifyJobError:      notifyJobError,
// 		NotifyTaskStarted:   notifyTaskStarted,
// 		NotifyTaskCompleted: notifyTaskCompleted,
// 		NotifyTaskError:     notifyTaskError,
// 	}, nil
// }
