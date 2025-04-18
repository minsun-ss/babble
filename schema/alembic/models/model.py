import sqlalchemy as sa
from sqlalchemy.dialects import mysql
from sqlalchemy.orm import declarative_base

Base = declarative_base()


class User(Base):
    __tablename__ = "users"

    id = sa.Column(sa.BigInteger, nullable=False, primary_key=True, autoincrement=True)
    username = sa.Column(sa.String(50), nullable=False, unique=True)
    iat = sa.Column(sa.Integer, nullable=False)
    last_updated_dt = sa.Column(
        mysql.TIMESTAMP,
        nullable=False,
        server_default=sa.text("current_timestamp() ON UPDATE current_timestamp()"),
    )


class UserAccess(Base):
    __tablename__ = "user_access"

    id = sa.Column(sa.BigInteger, nullable=False, primary_key=True, autoincrement=True)
    username = sa.Column(sa.String(50), sa.ForeignKey("users.username", ondelete="CASCADE"), nullable=False)
    project_name = sa.Column(sa.String(50), sa.ForeignKey("projects.project_name", ondelete="CASCADE"), nullable=False)
    last_updated_dt = sa.Column(
        mysql.TIMESTAMP,
        nullable=False,
        server_default=sa.text("current_timestamp() ON UPDATE current_timestamp()"),
    )

    __table_args__ = (
        sa.UniqueConstraint("username", "project_name", name="uix_username_projectname"),
        sa.Index("ix_username", username),
        sa.Index("ix_project_name", project_name),
        sa.Index("ix_last_updated_dt", last_updated_dt),
    )


class Projects(Base):
    __tablename__ = "projects"

    project_name = sa.Column(sa.String(50), nullable=False, primary_key=True)
    email = sa.Column(sa.String(50), nullable=True)
    last_updated_dt = sa.Column(
        mysql.TIMESTAMP,
        nullable=False,
        server_default=sa.text("current_timestamp() ON UPDATE current_timestamp()"),
    )

    __table_args__ = (
        sa.Index("ix_project_name", project_name),
        sa.Index("ix_last_updated_dt", last_updated_dt),
    )


class Docs(Base):
    __tablename__ = "docs"

    id = sa.Column(sa.BigInteger, nullable=False, autoincrement=True, primary_key=True)
    name = sa.Column(sa.String(50), nullable=False, unique=True)
    description = sa.Column(sa.String(50), nullable=True)
    is_visible = sa.Column(mysql.TINYINT, nullable=False, server_default=sa.text("1"))
    project_name = sa.Column(
        sa.String(50), sa.ForeignKey("projects.project_name", ondelete="CASCADE"), server_default=sa.text('"Other"')
    )
    last_updated_dt = sa.Column(
        mysql.TIMESTAMP,
        nullable=False,
        server_default=sa.text("current_timestamp() ON UPDATE current_timestamp()"),
    )

    __table_args__ = (
        sa.Index("ix_name", name),
        sa.Index("ix_project_name", project_name),
        sa.Index("ix_last_updated_dt", last_updated_dt),
    )


class DocHistory(Base):
    __tablename__ = "doc_history"

    id = sa.Column(sa.BigInteger, autoincrement=True, nullable=False, primary_key=True)
    name = sa.Column(sa.String(50), sa.ForeignKey("docs.name", ondelete="SET NULL"), nullable=True)
    version_major = sa.Column(sa.String(10), nullable=False)
    version_minor = sa.Column(sa.String(10), nullable=False)
    version_patch = sa.Column(sa.String(50), nullable=False)
    html = sa.Column(mysql.LONGBLOB, nullable=False)
    last_updated_dt = sa.Column(
        mysql.TIMESTAMP,
        nullable=False,
        server_default=sa.text("current_timestamp() ON UPDATE current_timestamp()"),
    )

    __table_args__ = (
        sa.Index("ix_version_major", version_major),
        sa.Index("ix_version_minor", version_minor),
        sa.Index("ix_version_patch", version_patch),
        sa.Index("ix_name", name),
        sa.UniqueConstraint(name, version_major, version_minor, version_patch, name="name"),
    )
