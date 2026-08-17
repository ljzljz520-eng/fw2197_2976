package main

import (
	"flag"
	"fmt"
	"os"

	"skillsassessment/domain"
	"skillsassessment/repository"
	"skillsassessment/service"
	"skillsassessment/storage"
)

func main() {
	path := "skills.db"
	if value := os.Getenv("SKILLS_DB"); value != "" {
		path = value
	}
	db, err := storage.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer db.Close()
	projects := repository.NewProjectRepository(db)
	learners := repository.NewLearnerRepository(projects)
	svc := service.NewProjectService(projects, learners)
	if err := execute(os.Args[1:], svc); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func execute(args []string, svc *service.ProjectService) error {
	if len(args) == 0 || args[0] == "help" {
		fmt.Println("skills assessment commands: register, publish, archive, list")
		return nil
	}
	switch args[0] {
	case "register":
		return registerCommand(args[1:], svc)
	case "publish":
		return transitionCommand(args[1:], svc, true)
	case "archive":
		return transitionCommand(args[1:], svc, false)
	case "list":
		projects, err := svc.ListProjects()
		if err != nil {
			return err
		}
		for _, project := range projects {
			fmt.Printf("%s\t%s\t%s\t%d learners\n", project.ID, project.Name, project.Status, len(project.Learners))
		}
		return nil
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func registerCommand(args []string, svc *service.ProjectService) error {
	flags := flag.NewFlagSet("register", flag.ContinueOnError)
	id := flags.String("id", "", "project id")
	name := flags.String("name", "", "project name")
	creator := flags.String("creator", "", "creator")
	if err := flags.Parse(args); err != nil {
		return err
	}
	project, err := svc.RegisterProject(*id, *name, *creator, defaultCertification())
	if err != nil {
		return err
	}
	fmt.Printf("registered %s\n", project.ID)
	return nil
}

func transitionCommand(args []string, svc *service.ProjectService, publish bool) error {
	flags := flag.NewFlagSet("transition", flag.ContinueOnError)
	id := flags.String("id", "", "project id")
	reviewer := flags.String("reviewer", "", "reviewer")
	if err := flags.Parse(args); err != nil {
		return err
	}
	var projectID string
	var err error
	if publish {
		result, publishErr := svc.PublishProject(*id, *reviewer)
		projectID = result.ID
		err = publishErr
	} else {
		result, archiveErr := svc.ArchiveProject(*id, *reviewer)
		projectID = result.ID
		err = archiveErr
	}
	if err != nil {
		return err
	}
	fmt.Printf("transitioned %s\n", projectID)
	return nil
}

func defaultCertification() domain.Certification {
	return domain.Certification{Kind: domain.CertificationGrade, Grade: "PASS"}
}
