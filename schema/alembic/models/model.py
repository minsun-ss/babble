import sqlalchemy as sa
from sqlalchemy.dialects import mysql
from sqlalchemy.orm import declarative_base

Base = declarative_base()


class Docs(Base):
    __tablename__ = "docs"

    name = sa.Column(sa.String(50), nullable=False, primary_key=True)
    description = sa.Column(sa.String(50), nullable=True)
    hidden = sa.Column(mysql.TINYINT, nullable=False, server_default=sa.text("0"))
    last_updated_dt = sa.Column(
        mysql.TIMESTAMP,
        nullable=False,
        server_default=sa.text("current_timestamp() ON UPDATE current_timestamp()"),
    )

    __table_args__ = (sa.Index("ix_last_updated_dt", last_updated_dt),)


class DocHistory(Base):
    __tablename__ = "doc_history"

    id = sa.Column(sa.BigInteger, autoincrement=True, nullable=False, primary_key=True)
    name = sa.Column(sa.String(50), nullable=False)
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
        sa.Index("idx_version_major", version_major),
        sa.Index("idx_version_minor", version_minor),
        sa.Index("idx_version_patch", version_patch),
        sa.Index("idx_name", name),
        sa.UniqueConstraint(name, version_major, version_minor, version_patch, name="name"),
    )
