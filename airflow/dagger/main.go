// A generated module for Airflow functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/airflow/internal/dagger"
)

type Airflow struct {
	// +private
	Ctr *dagger.Container
	// +private
	Redis *dagger.Redis
	// +private
	Postgres *dagger.Postgres
}

type serviceToBind struct {
	Alias string
	Svc   *dagger.Service
}

func New(
	ctx context.Context,
	// The image name to use
	// +default="apache/airflow"
	image string,
	// The version of apache airflow to sue
	// +default="2.9.3"
	version string,
	//
	// +optional
	dags *dagger.Directory,
	//
	// +optional
	config *dagger.Directory,
	//
	// +optional
	plugins *dagger.Directory,
	//
	// +optional
	requirements string,
	//
	// +optional
	// +default="default"
	databaseCacheName string,
) (*Airflow, error) {
	globalUser, globalPassword := "airflow", "airflow"
	commonEnvs := map[string]string{
		"AIRFLOW__CORE__EXECUTOR":                    "CeleryExecutor",
		"AIRFLOW__DATABASE__SQL_ALCHEMY_CONN":        "postgresql+psycopg2://" + globalUser + ":" + globalPassword + "@postgres/airflow",
		"AIRFLOW__CELERY__RESULT_BACKEND":            "db+postgresql://" + globalUser + ":" + globalPassword + "@postgres/airflow",
		"AIRFLOW__CELERY__BROKER_URL":                "redis://:@redis:6379/0",
		"AIRFLOW__CORE__FERNET_KEY":                  "",
		"AIRFLOW__CORE__DAGS_ARE_PAUSED_AT_CREATION": "true",
		"AIRFLOW__CORE__LOAD_EXAMPLES":               "true",
		"AIRFLOW__API__AUTH_BACKENDS":                "airflow.api.auth.backend.basic_auth,airflow.api.auth.backend.session",
		"AIRFLOW__SCHEDULER__ENABLE_HEALTH_CHECK":    "true",
	}

	if requirements != "" {
		commonEnvs["_PIP_ADDITIONAL_REQUIREMENTS"] = requirements
	}

	commonCtr, err := dag.
		Container().
		From(image + ":" + version).
		With(func(ctr *dagger.Container) *dagger.Container {
			for k, v := range commonEnvs {
				ctr = ctr.WithEnvVariable(k, v)
			}

			if dags != nil {
				ctr = ctr.WithDirectory("/opt/airflow/dags", dags)
			}

			if config != nil {
				ctr = ctr.WithDirectory("/opt/airflow/config", config)
			}

			if plugins != nil {
				ctr = ctr.WithDirectory("/opt/airflow/plugins", plugins)
			}

			return ctr
		}).
		Sync(ctx)
	if err != nil {
		return nil, err
	}

	//svcs := []serviceToBind{
	//	{
	//		Alias: "postgres",
	//		Svc: dag.Postgres(dagger.PostgresOpts{
	//			User:     dag.SetSecret("airflow-db-user", globalUser),
	//			Password: dag.SetSecret("airflow-db-password", globalPassword),
	//			Version:  "13",
	//			//DataVolume: dag.CacheVolume("airflow-database-data-" + databaseCacheName),
	//		}).Service(),
	//	},
	//	{
	//		Alias: "redis",
	//		Svc: dag.Redis(dagger.RedisOpts{
	//			Version: "7.4.0",
	//			Port:    6379,
	//		}).Server().AsService(),
	//	},
	//}

	return &Airflow{
		Ctr: commonCtr,
		//	Ctr: commonCtr.
		//		With(func(ctr *dagger.Container) *dagger.Container {
		//			return bindDependencies(ctr, svcs)
		//		}).
		//		WithServiceBinding(
		//			"airflow-scheduler",
		//			commonCtr.
		//				WithEnvVariable("_AIRFLOW_DB_MIGRATE", "true").
		//				WithEnvVariable("_AIRFLOW_WWW_USER_CREATE", "create").
		//				WithEnvVariable("_AIRFLOW_WWW_USER_USERNAME", globalUser).
		//				WithEnvVariable("_AIRFLOW_WWW_USER_PASSWORD", globalPassword).
		//				WithExposedPort(8974, dagger.ContainerWithExposedPortOpts{Protocol: dagger.Tcp}).
		//				With(func(ctr *dagger.Container) *dagger.Container {
		//					return bindDependencies(ctr, svcs)
		//				}).
		//				WithExec([]string{"scheduler"}, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		//				AsService(),
		//		).
		//		WithServiceBinding(
		//			"airflow-triggerer",
		//			commonCtr.
		//				With(func(ctr *dagger.Container) *dagger.Container {
		//					return bindDependencies(ctr, svcs)
		//				}).
		//				WithExec([]string{"triggerer"}, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		//				AsService(),
		//		).
		//		WithServiceBinding(
		//			"airflow-worker",
		//		commonCtr.
		//				WithEnvVariable("DUMB_INIT_SETSID", "0").
		//				With(func(ctr *dagger.Container) *dagger.Container {
		//					return bindDependencies(ctr, svcs)
		//				}).
		//				WithExec([]string{"celery", "worker"}, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		//				AsService(),
		//		).
		//		WithExec([]string{"webserver"}, dagger.ContainerWithExecOpts{UseEntrypoint: true}),
	}, nil
}

func (a *Airflow) Serve(ctx context.Context) *dagger.Service {
	svcs := []serviceToBind{
		{
			Alias: "postgres",
			Svc: dag.Postgres(dagger.PostgresOpts{
				User:     dag.SetSecret("airflow-db-user", "airflow"),
				Password: dag.SetSecret("airflow-db-password", "airflow"),
				Version:  "13",
				//DataVolume: dag.CacheVolume("airflow-database-data-" + databaseCacheName),
			}).Service(),
		},
		{
			Alias: "redis",
			Svc: dag.Redis(dagger.RedisOpts{
				Version: "7.4.0",
				Port:    6379,
			}).Server().AsService(),
		},
	}

	return a.
		Ctr.
		With(func(ctr *dagger.Container) *dagger.Container {
			return bindDependencies(ctr, svcs)
		}).
		WithServiceBinding(
			"airflow-scheduler",
			a.Ctr.
				WithEnvVariable("_AIRFLOW_DB_MIGRATE", "true").
				WithEnvVariable("_AIRFLOW_WWW_USER_CREATE", "create").
				WithEnvVariable("_AIRFLOW_WWW_USER_USERNAME", "airflow").
				WithEnvVariable("_AIRFLOW_WWW_USER_PASSWORD", "airflow").
				WithExposedPort(8974, dagger.ContainerWithExposedPortOpts{Protocol: dagger.Tcp}).
				With(func(ctr *dagger.Container) *dagger.Container {
					return bindDependencies(ctr, svcs)
				}).
				WithExec([]string{"scheduler"}, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
				AsService(),
		).
		WithServiceBinding(
			"airflow-worker",
			a.Ctr.
				WithEnvVariable("DUMB_INIT_SETSID", "0").
				With(func(ctr *dagger.Container) *dagger.Container {
					return bindDependencies(ctr, svcs)
				}).
				WithExec([]string{"celery", "worker"}, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
				AsService(),
		).
		WithExec([]string{"webserver"}, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		AsService()

	return a.Ctr.AsService()
}

func bindDependencies(ctr *dagger.Container, svcs []serviceToBind) *dagger.Container {
	for _, svc := range svcs {
		ctr = ctr.WithServiceBinding(svc.Alias, svc.Svc)
	}

	return ctr
}
