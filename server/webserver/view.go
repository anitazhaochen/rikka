package webserver

import (
	"net/http"

	"github.com/7sDream/rikka/api"
	"github.com/7sDream/rikka/common/util"
	"github.com/7sDream/rikka/plugins"
)

func viewHandleFunc(w http.ResponseWriter, r *http.Request) {
	taskID := util.GetTaskIDByRequest(r)
	context.TaskID = taskID

	l.Info("Recieve a view request of task", taskID, "from ip", r.RemoteAddr)

	l.Debug("Send a url request of task", taskID, "to plugin manager")

	var pURL *api.URL
	var err error
	if pURL, err = plugins.GetURL(taskID, r, nil); err != nil {
		// state is not finished or error when get url, use view.html
		templateFilePath := viewTemplateFilePath
		l.Warn("Can't get url of task", taskID, ":", err)
		l.Warn("Render template", viewTemplateFileName)
		err = util.RenderTemplate(templateFilePath, w, context)
		if util.ErrHandle(w, err) {
			// RenderTemplate error
			l.Error("Erro when render template", viewTemplateFileName, ":", err)
		} else {
			// successfully
			l.Info("Render template", viewTemplateFileName, "successfully")
		}
		return
	}

	// state is finished, use viewFinish.html
	l.Debug("Recieve url of task", taskID, ":", pURL.URL)
	templateFilePath := finishedViewTemplateFilePath
	context.URL = pURL.URL
	err = util.RenderTemplate(templateFilePath, w, context)
	if util.ErrHandle(w, err) {
		// RenderTemplate error
		l.Error("Error happened when render template", finishedViewTemplateFileName, ":", err)
	} else {
		// successfully
		l.Info("Render template", finishedViewTemplateFileName, "successfully")
	}
}

// ViewHandler handle requset ask for image view page(${ViewPath}<taskID>), use templates/view.html
// Only accept GET Method
func viewHandleGenerator() http.HandlerFunc {
	return util.RequestFilter(
		"", "GET", l,
		util.DisableListDir(l, viewHandleFunc),
	)
}
